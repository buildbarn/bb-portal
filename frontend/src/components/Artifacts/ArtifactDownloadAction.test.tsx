import { render, screen } from "@testing-library/react";
import { describe, expect, test } from "vitest";
import ArtifactDownloadAction from "./ArtifactDownloadAction";

describe("ArtifactDownloadAction", () => {
  test("renders Download button for local /api/v1/servefile URL", () => {
    render(
      <ArtifactDownloadAction
        downloadUrl="/api/v1/servefile/abc123/myfile.tar.gz"
        uri="bytestream://build.example.com/blobs/abc123/456"
        fileName="myfile.tar.gz"
      />,
    );

    const link = screen.getByRole("link", { name: /download/i });
    expect(link).toBeInTheDocument();
    expect(link).toHaveAttribute("download", "myfile.tar.gz");
    expect(link).toHaveAttribute("target", "_self");
    expect(link).not.toHaveAttribute("rel");
  });

  test("renders External button for https:// URL", () => {
    render(
      <ArtifactDownloadAction
        downloadUrl="https://storage.example.com/artifacts/myfile.tar.gz"
        uri="https://storage.example.com/artifacts/myfile.tar.gz"
        fileName="myfile.tar.gz"
      />,
    );

    const link = screen.getByRole("link", { name: /external/i });
    expect(link).toBeInTheDocument();
    expect(link).toHaveAttribute("target", "_blank");
    expect(link).toHaveAttribute("rel", "noopener noreferrer");
    expect(link).not.toHaveAttribute("download");
  });

  test("renders External button for http:// URL", () => {
    render(
      <ArtifactDownloadAction
        downloadUrl="http://storage.example.com/artifacts/myfile.tar.gz"
        uri="http://storage.example.com/artifacts/myfile.tar.gz"
        fileName="myfile.tar.gz"
      />,
    );

    const link = screen.getByRole("link", { name: /external/i });
    expect(link).toBeInTheDocument();
    expect(link).toHaveAttribute("target", "_blank");
    expect(link).toHaveAttribute("rel", "noopener noreferrer");
    expect(link).not.toHaveAttribute("download");
  });

  test("renders MinusOutlined icon with tooltip when downloadUrl is null and uri is provided", () => {
    render(
      <ArtifactDownloadAction
        downloadUrl={null}
        uri="bytestream://build.example.com/blobs/abc123/456"
        fileName="myfile.tar.gz"
      />,
    );

    // No button/link should be rendered
    expect(screen.queryByRole("link")).not.toBeInTheDocument();
    // The Ant Design Tooltip wraps its children; the icon should be present
    // via aria or role=img. We verify no Download/External text appears.
    expect(screen.queryByText(/download/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/external/i)).not.toBeInTheDocument();
  });

  test("renders MinusOutlined icon with 'No URI' tooltip when both downloadUrl and uri are null", () => {
    render(
      <ArtifactDownloadAction
        downloadUrl={null}
        uri={null}
        fileName="myfile.tar.gz"
      />,
    );

    expect(screen.queryByRole("link")).not.toBeInTheDocument();
    expect(screen.queryByText(/download/i)).not.toBeInTheDocument();
    expect(screen.queryByText(/external/i)).not.toBeInTheDocument();
  });

  test("renders MinusOutlined icon when downloadUrl is undefined", () => {
    render(
      <ArtifactDownloadAction
        downloadUrl={undefined}
        uri={undefined}
        fileName="myfile.tar.gz"
      />,
    );

    expect(screen.queryByRole("link")).not.toBeInTheDocument();
  });
});
