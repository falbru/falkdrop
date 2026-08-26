import { Item, ItemMedia, ItemContent, ItemTitle } from "@/components/ui/item";
import { Spinner } from "@/components/ui/spinner";

const PendingResourceItem = () => {
  return (
    <Item variant="outline">
      <ItemMedia variant="icon">
        <Spinner />
      </ItemMedia>
      <ItemContent>
        <ItemTitle>Uploading...</ItemTitle>
      </ItemContent>
    </Item>
  );
};

export default PendingResourceItem;
