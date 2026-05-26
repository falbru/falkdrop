import { Link, useParams } from "react-router";
import Button from "../components/ui/Button";
import Card from "../components/ui/Card";
import useDrop from "../features/drop/hooks/useGetDrop";

const GetDropPage = () => {
  const { dropId } = useParams();
  const drop = useDrop(dropId);

  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-4">
      <div className="w-full max-w-2xl">
        <h1 className="text-4xl font-bold mb-8 text-center text-neutral-100 uppercase">
          {drop.data?.id}
        </h1>

        {drop.isLoading && (
          <Card className="text-center py-8">
            <p className="text-neutral-400">Loading...</p>
          </Card>
        )}

        {drop.error && (
          <Card className="text-center py-8">
            <p className="text-neutral-400">Drop not found</p>
            <Link to="/" className="inline-block mt-4">
              <Button variant="secondary">Go Home</Button>
            </Link>
          </Card>
        )}

        {drop.data && (
          <Card>
            <h2 className="text-sm font-medium text-neutral-400 mb-3">Files</h2>
            <ul className="space-y-2">
              {drop.data.resources.length === 0 ? (
                <li className="text-sm text-neutral-500">
                  No files in this drop
                </li>
              ) : (
                drop.data.resources.map((res) => (
                  <li key={res.id} className="flex items-center gap-3">
                    <a
                      href={res.downloadURL}
                      className="flex-1 truncate text-sm text-neutral-300 hover:text-neutral-100 underline decoration-neutral-600 underline-offset-2"
                    >
                      {res.id}
                    </a>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() =>
                        navigator.clipboard.writeText(res.downloadURL)
                      }
                    >
                      Copy
                    </Button>
                    <a href={res.downloadURL}>
                      <Button variant="primary" size="sm">
                        Download
                      </Button>
                    </a>
                  </li>
                ))
              )}
            </ul>
          </Card>
        )}

        <div className="mt-6 text-center">
          <Link to="/">
            <Button variant="ghost">Back to Home</Button>
          </Link>
        </div>
      </div>
    </div>
  );
};

export default GetDropPage;
