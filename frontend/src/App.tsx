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
import ParentView, { loader as indexLoader } from './routes/index.tsx';

//enum LoadedState { loading, loaded, failed }

const StationData = () => {
  return <p>Heisann!</p>;
}

function App() {
  const router = createBrowserRouter([
    {
      path: "/",
      element: <ParentView />,
      loader: indexLoader,
      children: [
        {
          path: "/stations",
          element: <p>Lastar inn...</p>,
          loader: async ({request}: LoaderFunctionArgs) => {
          }
        },
        {
          path: "/s/:stationId",
          element: <StationData />,
          loader: (...args) => {
            return args;
          }
        },
        {
          path: "/s/:stationId/from/:fromYear/to/:toYear",
          element: <StationData />,
          loader: (...args) => {
            return args;
          }
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
