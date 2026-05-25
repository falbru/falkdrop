import { useParams } from "react-router";
import useDrop from "../features/drop/hooks/useGetDrop";

const GetDropPage = () => {
  const { dropId } = useParams();

  const drop = useDrop(dropId);

  return (
    <div>
      <h1>{drop.data?.id}</h1>
      <ul>
        {drop.data?.resources.map((res) => (
          <li key={res.id}>
            <a href={res.downloadURL}>DOWNLOAD</a> {res.id}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default GetDropPage;
