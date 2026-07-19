import { useParams, useNavigate } from "react-router";
import { Button } from "../components/ui/button";
import { ItemGroup } from "../components/ui/item";
import ResourceItem from "../features/drop/components/ResourceItem";
import useDrop from "../features/drop/hooks/useGetDrop";
import { Download } from "lucide-react";
import { useEffect } from "react";
import { toast } from "sonner";

const GetDropPage = () => {
  const { dropId } = useParams();
  const drop = useDrop(dropId);
  const navigate = useNavigate();

  useEffect(() => {
    if (drop.error) {
      const errorMessage =
        drop.error instanceof Error
          ? drop.error.message
          : "Error occurred while loading drop";
      toast.error(errorMessage);
      navigate("/");
    }
  }, [drop.error, navigate]);

  return (
    <div className="flex flex-col gap-6">
      <h1 className="font-(family-name:--font-title) text-8xl font-bold text-center text-text uppercase">
        {drop.data?.id}
      </h1>

      {drop.data && (
        <ItemGroup title="Files">
          {drop.data.resources.map((res, index) => (
            <ResourceItem
              key={res.id}
              name={res.name ?? res.id}
              variant={
                index < drop.data.resources.length - 1 ? "default" : "muted"
              }
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
