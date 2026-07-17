import { Link, useParams } from "react-router";
import Button from "../components/ui/Button";
import Card from "../components/ui/Card";
import ItemGroup from "../components/ui/ItemGroup";
import ResourceItem from "../features/drop/components/ResourceItem";
import useDrop from "../features/drop/hooks/useGetDrop";
import { Download } from "lucide-react";

const GetDropPage = () => {
  const { dropId } = useParams();
  const drop = useDrop(dropId);

  return (
    <div className="flex flex-col gap-6">
      <h1 className="font-(family-name:--font-title) text-8xl font-bold text-center text-text uppercase">
        {drop.data?.id}
      </h1>

      {drop.error ? (
        <Card className="text-center py-8">
          <p className="text-text-secondary">Drop not found</p>
          <Link to="/" className="inline-block mt-4">
            <Button variant="secondary">Go Home</Button>
          </Link>
        </Card>
      ) : (
        <></>
      )}

      {drop.data && (
        <ItemGroup title="Files" empty="No files in this drop">
          {drop.data.resources.map((res, index) => (
            <ResourceItem
              key={res.id}
              name={res.name ?? res.id}
              showBorder={index < drop.data.resources.length - 1}
            >
              <a href={res.downloadURL}>
                <Button variant="ghost">
                  <Download />
                </Button>
              </a>
            </ResourceItem>
          ))}
        </ItemGroup>
      )}
    </div>
  );
};

export default GetDropPage;
