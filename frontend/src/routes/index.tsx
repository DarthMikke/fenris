
import {
  Link,
  useOutlet,
  useLoaderData,
  LoaderFunctionArgs
} from 'react-router-dom';

import stations from '../mock_data/stations.json';
console.debug(stations);

type StationsListProps = {
  stations: Station[]
} & React.PropsWithChildren;

const Stations = ({stations}: StationsListProps) => {
  return <ul>
    {stations.map(station => <li key={station.id}>
      <Link to={`/s/${station.id}`}>{station.name}</Link>
    </li>)}
  </ul>;
}

export type Station = {
  "name": string,
  "municipality": string,
  "county": string,
  "id": string,
  "elevation": number,
  "latitude": number,
  "longitude": number,
  "validFrom": string|null,
  "validTo": string|null,
  "wmo": number,
  "wigos": string,
  "owner": string
};

type StationsLoaderData = {
  stations: Station[],
};

export function loader(_: LoaderFunctionArgs) {
  return stations;
}

export default () => {
  const mainView = useOutlet();
  const stations = useLoaderData() as StationsLoaderData;

  return <>
  {mainView || ( stations ? <Stations stations={stations.stations} /> : <p>Loading...</p> )}
  </>;
}
