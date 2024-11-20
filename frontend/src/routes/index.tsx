
import {
  Link,
  useOutlet,
  useLoaderData,
  LoaderFunctionArgs
} from 'react-router-dom';

const Stations = ({stations}: {stations: any[]} & React.PropsWithChildren) => {
  return <ul>
    {stations.map(station => <li key={station.id}>
      <Link to={`/s/${station.id}`}>{station.name}</Link>
    </li>)}
  </ul>
}

export async function loader({ params, request }: LoaderFunctionArgs) {
  return {
    "stations": [
      {
        "name": "Lillehammer - Sætherengen",
        "municipality": "Lillehammer",
        "county": "Innlandet",
        "id": "SN12680",
        "elevation": 240,
        "latitude": 61.0917,
        "longitude": 10.4762,
        "availableFrom": "1982-12-01",
        "availableTo": null,
        "wmo": 1378,
        "wigos": "0-20000-0-01378",
        "owner": "Met.no",
      }
    ]
  }
}

export default () => {
  const mainView = useOutlet();
  const stations = useLoaderData();

  return <>
  {mainView || ( stations ? <Stations stations={stations.stations} /> : <p>Loading...</p> )}
  </>;
}
