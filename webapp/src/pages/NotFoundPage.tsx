import Title from "@/components/shared/Title";
import { LinkButton } from "@/components/ui/button";

const NotFoundPage = () => {
  return (
    <div className="flex flex-col items-center text-center">
      <Title>404 Not Found</Title>
      <h1 className="font-(family-name:--font-title) text-6xl font-bold mb-4 uppercase">
        404
      </h1>
      <p className="text-lg mb-2">Page not found</p>
      <p className="text-muted-foreground mb-8 max-w-md">
        The page you're looking for doesn't exist.
      </p>
      <LinkButton href="/">Go Home</LinkButton>
    </div>
  );
};

export default NotFoundPage;
