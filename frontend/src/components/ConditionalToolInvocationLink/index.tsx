import { useQuery } from "@apollo/client/react";
import { Link } from "@tanstack/react-router";
import { CHECK_IF_INVOCATION_EXISTS } from "./graphql";

interface Props {
  toolInvocationID: string;
}

const ConditionalToolInvocationLink: React.FC<Props> = ({
  toolInvocationID,
}) => {
  const { data } = useQuery(CHECK_IF_INVOCATION_EXISTS, {
    variables: { invocationID: toolInvocationID },
    fetchPolicy: "cache-and-network",
  });

  if (data?.getBazelInvocation !== undefined)
    return (
      <Link
        to={`/bazel-invocations/$invocationID`}
        params={{ invocationID: toolInvocationID }}
      >
        {toolInvocationID}
      </Link>
    );

  return <>{toolInvocationID}</>;
};

export default ConditionalToolInvocationLink;
