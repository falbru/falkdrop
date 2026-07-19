import { LinkButton } from "../../../components/ui/button";

export default function DropNotFound() {
  return (
    <div className="flex flex-col items-center text-center">
      <h1 className="font-(family-name:--font-title) text-6xl font-bold mb-4 uppercase">
        404
      </h1>
      <p className="text-lg mb-2">Drop not found</p>
      <p className="text-muted-foreground mb-8 max-w-md">
        The drop you're looking for doesn't exist or may have expired.
      </p>
      <LinkButton href="/">Go Home</LinkButton>
    </div>
  );
}
