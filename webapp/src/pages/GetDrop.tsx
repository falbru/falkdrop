import { useParams } from "react-router";
import { Button } from "../components/ui/button";
import { ItemGroup } from "../components/ui/item";
import ResourceItem from "../features/drop/components/ResourceItem";
import DropNotFound from "../features/drop/components/DropNotFound";
import useDrop from "../features/drop/hooks/useGetDrop";
import { Copy, Download, QrCode as QrCodeIcon } from "lucide-react";
import { Spinner } from "../components/ui/spinner";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from "@/components/ui/input-group";
import { toast } from "sonner";
import QRCode from "react-qr-code";
import {
  Dialog,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../components/ui/dialog";

interface LinkInputProps {
  link: string;
  onCopy: () => void;
}

interface QrCodeDialogProps {
  link: string;
}

const LinkInput = ({ link, onCopy }: LinkInputProps) => (
  <InputGroup>
    <InputGroupInput value={link} disabled />
    <InputGroupAddon align="inline-end">
      <InputGroupButton variant="ghost" size="icon-sm" onClick={onCopy}>
        <Copy />
      </InputGroupButton>
    </InputGroupAddon>
  </InputGroup>
);

const QrCodeDialog = ({ link }: QrCodeDialogProps) => (
  <DialogTrigger>
    <Button variant="secondary" size="icon">
      <QrCodeIcon />
    </Button>
    <Dialog>
      <DialogHeader>
        <DialogTitle>Share with QR Code</DialogTitle>
      </DialogHeader>
      <div className="flex justify-center p-4">
        <div className="p-4 rounded-lg bg-white">
          <QRCode value={link} />
        </div>
      </div>
    </Dialog>
  </DialogTrigger>
);

const GetDropPage = () => {
  const { dropId } = useParams();
  const drop = useDrop(dropId);
  const link = `${window.origin}/${dropId ?? ""}`;

  const handleCopyLink = async () => {
    try {
      await navigator.clipboard.writeText(link);
      toast.success("Link copied to clipboard!");
    } catch (error) {
      toast.error("Failed to copy link");
      console.error("Failed to copy link", error);
    }
  };

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

      <div className="flex gap-1">
        <LinkInput
          link={link}
          onCopy={() => {
            void handleCopyLink();
          }}
        />
        <QrCodeDialog link={link} />
      </div>

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
