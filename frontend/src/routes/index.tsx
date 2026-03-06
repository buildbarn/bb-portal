import { createFileRoute } from '@tanstack/react-router';
import { HomePage } from "@/components/pages/Home";
import { generatePageTitle } from '@/utils/generatePageTitle';

export const Route = createFileRoute('/')({
  component: HomePage,
  head: (_ctx) => ({meta: [{title: generatePageTitle(["Overview"])}]})
})
