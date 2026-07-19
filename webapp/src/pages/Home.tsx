import { useState } from "react";
import { useNavigate } from "react-router";
import { Button, LinkButton } from "../components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import { Plus } from "lucide-react";
import { useAuth } from "@/features/auth/contexts/AuthProvider";
import { useQueryClient } from "@tanstack/react-query";
import fetchDrop from "@/features/drop/api/fetchDrop";
import { toast } from "sonner";

const HomePage = () => {
  const [inputDropCode, setInputDropCode] = useState<string>("");
  const [isValidating, setIsValidating] = useState<boolean>(false);
  const navigate = useNavigate();

  const auth = useAuth();
  const queryClient = useQueryClient();

  const handleEnter = async () => {
    if (inputDropCode.length < 5) return;

    setIsValidating(true);

    try {
      await queryClient.fetchQuery({
        queryKey: ["drop", inputDropCode],
        queryFn: () => fetchDrop(inputDropCode),
      });

      navigate(`/${inputDropCode}`);
    } catch (error) {
      const errorMessage =
        error instanceof Error ? error.message : "Drop not found";

      toast.error(errorMessage);
    } finally {
      setIsValidating(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleEnter();
    }
  };

  return (
    <div className="flex flex-col items-center flex-grow justify-center">
      <h1 className="font-(family-name:--font-title) uppercase text-6xl font-bold mb-6 text-center">
        Falk
        <br />
        Drop
      </h1>

      <Card className="w-[480px] max-w-full">
        <CardContent className="flex flex-col gap-2">
          <Input
            onChange={(e) => setInputDropCode(e.target.value)}
            value={inputDropCode}
            placeholder="Drop Code"
            maxLength={5}
            className="uppercase text-center"
            onKeyDown={handleKeyDown}
          />

          <Button
            onClick={handleEnter}
            isDisabled={inputDropCode.length < 5 || isValidating}
            className="w-full"
          >
            Enter
          </Button>
        </CardContent>

        {auth && auth.isAuthenticated() && (
          <CardFooter className="flex justify-center">
            <LinkButton
              href="/create"
              variant="secondary"
              className="flex gap-1 items-center"
            >
              <Plus size={18} /> <span>Create Drop</span>
            </LinkButton>
          </CardFooter>
        )}
      </Card>
    </div>
  );
};

export default HomePage;
