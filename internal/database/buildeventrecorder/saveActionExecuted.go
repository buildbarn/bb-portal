package buildeventrecorder

import (
	"context"
	"reflect"

	bes "github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/buildeventstream/proto"
	"github.com/bazelbuild/bazel/src/main/protobuf"
	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocation"
	"github.com/buildbarn/bb-portal/ent/gen/ent/configuration"
	"github.com/buildbarn/bb-storage/pkg/util"
)

func getErrorCodeFromFailureDetail(failureDetail *protobuf.FailureDetail) string {
	if failureDetail == nil || failureDetail.Category == nil {
		return ""
	}
	detailValue := reflect.ValueOf(failureDetail.Category)
	if detailValue.Kind() == reflect.Ptr {
		detailValue = detailValue.Elem()
	}
	if detailValue.Kind() != reflect.Struct {
		return ""
	}

	for i := 0; i < detailValue.NumField(); i++ {
		fieldValue := detailValue.Field(i)
		if fieldValue.Kind() != reflect.Ptr {
			continue
		}

		method := fieldValue.MethodByName("GetCode")
		if !method.IsValid() {
			continue
		}

		result := method.Call(nil)
		if len(result) == 0 {
			continue
		}

		stringer := result[0].MethodByName("String")
		if !stringer.IsValid() {
			continue
		}
		return stringer.Call(nil)[0].String()
	}
	return ""
}

func (r *BuildEventRecorder) saveActionExecuted(ctx context.Context, tx *ent.Client, actionExecuted *bes.ActionExecuted, actionCompletedID *bes.BuildEventId_ActionCompletedId) error {
	if actionExecuted == nil || actionCompletedID == nil {
		return nil
	}
	if actionCompletedID.Label == "" {
		return nil
	}
	// We are only interested in failed actions. If this is changed, some of
	// the text in the frontend needs to be updated as well.
	if actionExecuted.Success {
		return nil
	}

	create := tx.Action.Create().
		SetBazelInvocationID(r.InvocationDbID).
		SetLabel(actionCompletedID.Label).
		SetSuccess(actionExecuted.Success).
		SetExitCode(actionExecuted.ExitCode).
		SetCommandLine(actionExecuted.CommandLine)

	if configID := actionCompletedID.Configuration.GetId(); configID != "" {
		// This results in a database query per ActionExecuted event. This is
		// acceptable since we only care about failed actions, which are
		// relatively rare. If we ever care about successful actions as well,
		// we should batch this work.
		configDbID, err := tx.Configuration.Query().
			Where(
				configuration.ConfigurationID(configID),
				configuration.HasBazelInvocationWith(bazelinvocation.ID(r.InvocationDbID)),
			).
			OnlyID(ctx)
		if err != nil {
			return util.StatusWrapf(err, "failed to query Configuration with ID %#v for ActionExecuted", configID)
		}
		create.SetConfigurationID(configDbID)
	}

	if actionExecuted.Type != "" {
		create.SetType(actionExecuted.Type)
	}
	if failureMessage := actionExecuted.GetFailureDetail().GetMessage(); failureMessage != "" {
		create.SetFailureMessage(failureMessage)
	}
	if failureCode := getErrorCodeFromFailureDetail(actionExecuted.GetFailureDetail()); failureCode != "" {
		create.SetFailureCode(failureCode)
	}
	if actionExecuted.StartTime != nil {
		create.SetStartTime(actionExecuted.EndTime.AsTime())
	}
	if actionExecuted.EndTime != nil {
		create.SetEndTime(actionExecuted.EndTime.AsTime())
	}
	if stdoutFile := actionExecuted.GetStdout(); stdoutFile != nil {
		if digest := getFileDigestFromBesFile(stdoutFile); digest != nil {
			create.SetStdoutHash(*digest)
		}
		if length := getFileSizeBytesFromBesFile(stdoutFile); length != nil {
			create.SetStdoutSizeBytes(*length)
		}
		if digestFunction := getFileDigestFunctionFromBesFile(stdoutFile); digestFunction != nil {
			create.SetStdoutHashFunction(*digestFunction)
		}
	}
	if stderrFile := actionExecuted.GetStderr(); stderrFile != nil {
		if digest := getFileDigestFromBesFile(stderrFile); digest != nil {
			create.SetStderrHash(*digest)
		}
		if length := getFileSizeBytesFromBesFile(stderrFile); length != nil {
			create.SetStderrSizeBytes(*length)
		}
		if digestFunction := getFileDigestFunctionFromBesFile(stderrFile); digestFunction != nil {
			create.SetStderrHashFunction(*digestFunction)
		}
	}
	err := create.Exec(ctx)
	if err != nil {
		return util.StatusWrap(err, "failed to save Action")
	}
	return nil
}
