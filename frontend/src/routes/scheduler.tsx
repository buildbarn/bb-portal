import { env } from '@/utils/env';
import { requireFeature } from '@/utils/featureGuard';
import { createFileRoute, Outlet } from '@tanstack/react-router';

export const Route = createFileRoute('/scheduler')({
  component: Outlet,
  beforeLoad: requireFeature(env.featureFlags?.scheduler),
})
