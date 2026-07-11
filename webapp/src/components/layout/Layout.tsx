import { Link, Outlet } from "react-router";
import { useAuth } from "../../features/auth/contexts/AuthProvider";

type LayoutProps = {
  showHeader?: boolean;
};

const Layout = (props: LayoutProps) => {
  const { showHeader = true } = props;
  const auth = useAuth();

  return (
    <div className="w-[680px] max-w-screen min-h-screen flex flex-col mx-auto">
      {showHeader && (
        <header className="py-4">
          <Link to="/">
            <h1 className="font-(family-name:--font-title) uppercase text-center text-3xl font-bold text-neutral-100 hover:text-neutral-300 transition-colors">
              FalkDrop
            </h1>
          </Link>
        </header>
      )}
      <main className="flex-grow p-8">
        <Outlet />
      </main>
      <footer className="py-4 text-center">
        {auth && auth.isAuthenticated() ? (
          <button
            onClick={() => auth.logout()}
            className="text-sm text-neutral-400 hover:text-neutral-200 transition-colors"
          >
            Log out
          </button>
        ) : (
          <Link
            to="/login"
            className="text-sm text-neutral-400 hover:text-neutral-200 transition-colors"
          >
            Log in
          </Link>
        )}
      </footer>
    </div>
  );
};

export default Layout;
