import { useQuery } from "@apollo/client";
import Link from "next/link";
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
      <Link href={`/bazel-invocations/${toolInvocationID}`}>
        {toolInvocationID}
      </Link>
    );

  return <>{toolInvocationID}</>;
};

export default ConditionalToolInvocationLink;
