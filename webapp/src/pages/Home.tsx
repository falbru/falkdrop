import { Input } from "react-aria-components";
import Button from "../components/ui/Button";

const HomePage = () => {
  return (
    <>
      <h1>FalkDrop</h1>
      <div className="flex flex-col w-[500px]">
        <Input />
        <Button>Enter</Button>
      </div>
    </>
  );
};

export default HomePage;
