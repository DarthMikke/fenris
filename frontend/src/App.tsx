import './App.scss';
import React from 'react';
import {
  Link,
  LoaderFunctionArgs,
  RouterProvider,
  createBrowserRouter,
  redirect,
  useOutlet,
  useLoaderData
} from 'react-router-dom';
import RootView, { loader as indexLoader } from './routes/index.tsx';
import StationView, { loader as stationLoader } from './routes/s/_stationId.tsx';
import StatsView, { loader as statsLoader } from './routes/s/_stationId/from/_fromYear/to/_toYear.tsx';

//enum LoadedState { loading, loaded, failed }

console.log()

function App() {
  const router = createBrowserRouter([
    {
      path: "/",
      element: <RootView />,
      loader: indexLoader,
      children: [
        {
          path: "s/:stationId",
          id: "station",
          element: <StationView />,
          loader: stationLoader,
          children: [
            {
              path: "from/:fromYear/to/:toYear",
              element: <StatsView />,
              loader: statsLoader,
            },
          ],
        },
        {
          path: "/test/",
          element: <p>Heisann!</p>,
        }
      ]
    }
  ]);

  return (
    <>
        <RouterProvider router={router} />
    </>
  )
}

export default App
