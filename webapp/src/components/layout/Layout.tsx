import { Link, Outlet } from "react-router";

const Layout = () => {
  return (
    <div className="w-[680px] max-w-screen min-h-screen flex flex-col mx-auto">
      <header className="py-4">
        <Link to="/">
          <h1 className="font-(family-name:--font-title) uppercase text-center text-3xl font-bold text-neutral-100 hover:text-neutral-300 transition-colors">
            FalkDrop
          </h1>
        </Link>
      </header>
      <main className="p-8">
        <Outlet />
      </main>
    </div>
  );
};

export default Layout;
