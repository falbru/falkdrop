import { useParams } from "react-router";
import { Button } from "../components/ui/button";
import { ItemGroup } from "../components/ui/item";
import ResourceItem from "../features/drop/components/ResourceItem";
import DropNotFound from "../features/drop/components/DropNotFound";
import useDrop from "../features/drop/hooks/useGetDrop";
import { Download } from "lucide-react";
import { Spinner } from "../components/ui/spinner";

const GetDropPage = () => {
  const { dropId } = useParams();
  const drop = useDrop(dropId);

  if (drop.isLoading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh]">
        <Spinner className="w-12 h-12" />
      </div>
    );
  }

  if (drop.error || !dropId) {
    return <DropNotFound />;
  }

  if (!drop.data) {
    return <DropNotFound />;
  }

  return (
    <div className="flex flex-col gap-6">
      <h1 className="font-(family-name:--font-title) text-8xl font-bold text-center text-text uppercase">
        {drop.data.id}
      </h1>

      <ItemGroup title="Files">
        {drop.data.resources.map((res) => (
          <ResourceItem key={res.id} name={res.name ?? res.id}>
            <a href={res.downloadURL}>
              <Button variant="ghost">
                <Download />
              </Button>
            </a>
          </ResourceItem>
        ))}
      </ItemGroup>
    </div>
  );
};

export default GetDropPage;
