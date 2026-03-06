import { generatePageTitle } from '@/utils/generatePageTitle';
import { OperationsPage } from '@/components/pages/Operations';
import { createFileRoute } from '@tanstack/react-router';
import { z } from 'zod';

const OperationsFilterSchema = z.object({
  "@type": z.string(),
  toolInvocationId: z.uuid().optional(),
  correlatedInvocationsId: z.uuid().optional(),
}).optional().refine(
  data => (data?.toolInvocationId && !data?.correlatedInvocationsId) || (!data?.toolInvocationId && data?.correlatedInvocationsId),
  "Either toolInvocationId or correlatedInvocationsId must be provided, but not both"
)

export type OperationsFilterParams = z.infer<typeof OperationsFilterSchema>

const OperationsSearchSchema = z.object({
  filter: OperationsFilterSchema
})

export type OperationsSearchParams = z.infer<typeof OperationsSearchSchema>

export const Route = createFileRoute('/operations/')({
  component: RouteComponent,
  validateSearch: (search) => OperationsSearchSchema.parse(search),
  head: (_ctx) => ({ meta: [{ title: generatePageTitle(["Operations"]) }] }),
})

function RouteComponent() {
  const { filter } = Route.useSearch();
  return <OperationsPage filter={filter} />
}
