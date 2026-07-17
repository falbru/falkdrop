import { useState } from "react";
import { Link, useNavigate } from "react-router";
import Button from "../components/ui/Button";
import Input from "../components/ui/Input";
import Card from "../components/ui/Card";
import { Plus } from "lucide-react";
import { useAuth } from "../features/auth/contexts/AuthProvider";

const HomePage = () => {
  const [inputDropCode, setInputDropCode] = useState<string>("");
  const navigate = useNavigate();

  const auth = useAuth();

  const handleEnter = () => {
    if (inputDropCode.length < 5) return;
    navigate(`/${inputDropCode}`);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleEnter();
    }
  };

  return (
    <div className="flex flex-col items-center flex-grow justify-center">
      <h1 className="font-(family-name:--font-title) uppercase text-6xl font-bold mb-6 text-center text-text">
        Falk
        <br />
        Drop
      </h1>

      <Card className="flex flex-col gap-3 w-[400px] max-w-full">
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
          variant="primary"
          isDisabled={inputDropCode.length < 5}
        >
          Enter
        </Button>

        {auth && auth.isAuthenticated() && (
          <>
            <div className="flex items-center gap-2 my-3">
              <div className="flex-grow h-[1px] bg-border" />
              <span className="text-sm text-text-muted uppercase">Or</span>
              <div className="flex-grow h-[1px] bg-border" />
            </div>

            <div className="flex justify-center">
              <Link to="/create">
                <Button variant="secondary" className="flex gap-1 items-center">
                  <Plus size={18} /> <span>Create Drop</span>
                </Button>
              </Link>
            </div>
          </>
        )}
      </Card>
    </div>
  );
};

export default HomePage;
