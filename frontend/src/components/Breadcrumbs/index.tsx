'use client';

import React, { useMemo } from 'react';
import { Breadcrumb, Typography } from 'antd';
import { BreadcrumbItemType } from 'antd/lib/breadcrumb/Breadcrumb';
import { HomeOutlined } from '@ant-design/icons';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import styles from '@/components/Breadcrumbs/index.module.css';

// Define a common separator reference
const separator: string = '/';

// Implements a custom breadcrumb item renderer because link of Breadcrumb targets # by default
function itemRender(item: BreadcrumbItemType, params: object, routes: BreadcrumbItemType[], paths: string[]) {
  // Determine whether the item is the last in the list of breadcrumbs to be rendered
  const lastItem = routes.findIndex(route => route.path === item.path) === routes.length - 1;

  // Return plain text for the last item since we are already on that page
  // Return a link to the item for all parent pages
  return lastItem || item.path === undefined ? (
    <Typography.Text type="secondary">{item.title}</Typography.Text>
  ) : (
    <Link href={item.path}>{item.title}</Link>
  );
}

// Allow pretty segment titles to be provided
// These are to override the default behavior for converting path segments to breadcrumbs
interface Props {
  segmentTitles?: string[];
}

// Compose a list of breadcrumb items for the current path
const Breadcrumbs: React.FC<Props> = ({ segmentTitles }) => {
  // Fetch the path name from Next
  const pathname: string = usePathname();

  // Generate the list of breadcrumb items
  // Memoize to avoid re-rending the breadcrumb every time the component is re-rendered
  const breadcrumbs: BreadcrumbItemType[] = useMemo(
    // Convert a path into a list of breadcrumb items to pass to an Ant Breadcrumb component
    () => {
      // Split the path by slashes and remove any empty segments
      // For example, convert "/foo/bar/baz" to ["foo", "bar", "baz"]
      const segments = pathname.split(separator).filter(segment => segment.length > 0);

      // Iterate over the list of segments to build a breadcrumb object for each
      const breadcrumbs = segments.map((segment, index) => {
        // Compose the path by joining all preceding segments
        const path = separator + segments.slice(0, index + 1).join(separator);

        // Use the segment as the title to be displayed in the browser
        // Convert to upper case to avoid implementing a title case converter
        const title = segmentTitles && segmentTitles.length > index ? segmentTitles[index] : segment.toUpperCase();

        // Return the path and title for the breadcrumb item
        return { path, title };
      });

      // Return the list of breadcrumbs items along with a way home
      return [{ path: separator, title: <HomeOutlined /> }, ...breadcrumbs];

      // Only update the breadcrumb item list when the path has changed
    },
    [pathname, segmentTitles],
  );

  // Compose and return the breadcrumb component
  // Use the custom breadcrumb item renderer because link of Breadcrumb targets # by default
  return (
    <div className={styles.breadcrumbs}>
      <Breadcrumb itemRender={itemRender} items={breadcrumbs} />
    </div>
  );
};

export default Breadcrumbs;
