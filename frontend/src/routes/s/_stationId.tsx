import {
  Link,
  LoaderFunctionArgs,
  useOutlet,
  useLoaderData,
  useNavigate
} from 'react-router-dom';
import { Station } from '../index';
import { useEffect } from 'react';
import { stations } from '../../mock_data/stations.json';

export const loader = async ({ params }: LoaderFunctionArgs) => {
  const response = await fetch(`/api/s/${params.stationId}`);
  const station = (await response.json()) as Station;

  return station;
};

export default () => {
  const station = useLoaderData() as Station;
  const statsView = useOutlet();
  const navigate = useNavigate();

  useEffect(() => {
    if (statsView === null) {
      navigate('from/1990/to/2020');
    }
  }, [statsView]);

  const availableFrom = station.validFrom;
  const availableTo = station.validTo;

  return <>
    <div className="flex-fill">
      <h1 className='col-12'>{station.name}</h1>
      <div className='row flex'>
      <div className="col-md-4">
        <p>Målestasjon i <strong>{station.municipality}, {station.county}</strong>.</p>
        <p>Data tilgjengeleg frå <strong>{availableFrom}</strong> til <strong>
        {availableTo || "no"}</strong>.</p>
      </div>
      {statsView}
      </div>
    </div>
  </>;
};
