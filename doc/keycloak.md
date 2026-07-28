# Setting up keycloak with FalkDrop

When having a running keycloak instance running, you have to do some manual set up to make it ready for use with falkdrop.

1. Create a new realm for falkdrop (used by `VITE_KEYCLOAK_REALM` and `KEYCLOAK_REALM_URL` env variables in the frontend and backend respectively).

2. Create the `falkdrop-api` client. Set the Valid Redirect URL and Web Origins to the backend URL (e.g., `http://localhost:8082/*` by default).

3. Create the `falkdrop-webapp` client. Set the Valid Redirect URL to the frontend URL (e.g., `http://localhost:5173`) and Web Origins to `*` (required for SPA mode without server-side rendering). Add `falkdrop-api` as its audience by going to `Client scopes -> falkdrop-webapp-dedicated -> Configure a new mapper -> Audience`.

4. Create a new Realm role called `drops:create`.

5. Create users as needed. Users who should be able to create drops should be assigned the Realm role `drops:create`.
