import { generatePageTitle } from '@/utils/generatePageTitle';
import { BrowserWelcomePage } from '@/components/pages/Browser/BrowserWelcome';
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/browser/')({
  component: RouteComponent,
  head: (_ctx) => ({meta: [{title: generatePageTitle(["Browser"])}]})
});

function RouteComponent() {
  return <BrowserWelcomePage />;
}
