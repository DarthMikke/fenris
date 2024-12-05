import { useState } from 'react';
import {
  Link,
  LoaderFunctionArgs,
  useLoaderData,
  useRouteLoaderData
} from 'react-router-dom';
import { Station } from '../../../../..';

export const loader = async ({params}: LoaderFunctionArgs) => {
  const resp = await fetch(`/api/s/${params.stationId}/from/${params.fromYear}/to/${params.toYear}`);
  const data = await resp.json()
  return {
    ...params,
    series: data,
    loaded: true,
  };
};

const months = ['Jan', 'Feb', 'Mar', 'Apr', 'Mai', 'Jun', 'Jul', 'Aug', 'Sep', 'Okt', 'Nov', 'Des'];

export default () => {
  const stats = useLoaderData() as {
    fromYear: string,
    toYear: string,
    series: {[id: string]: number[]},
    loaded: boolean
  };
  const [selectedFromYear, setSelectedFromYear] = useState(stats.fromYear);
  const [selectedToYear, setSelectedToYear] = useState(stats.toYear);
  const station = useRouteLoaderData("station") as Station;

  const table = stats.loaded
    ? <table className='table'>
      <thead>
        <tr>
          <th>Månad</th>
          <th>Gjennomsnittleg minimumstemperatur</th>
          <th>Gjennomsnittleg temperatur</th>
          <th>Gjennomsnittleg maksimumstemperatur</th>
        </tr>
      </thead>
      <tbody>
        {months.map((month, i) => {
          return <tr>
            <td>{month}</td>
            <td>{stats.series.avgmin[i]}</td>
            <td>{stats.series.average[i]}</td>
            <td>{stats.series.avgmax[i]}</td>
          </tr>
        })}
      </tbody>
    </table>
    : <p>Dataa dine er ikkje klare enno, prøv igjen snart.</p>;

  return <div className="col-md-8">
    <p>Frå: <input type='number' value={selectedFromYear}
      onChange={(e) => { setSelectedFromYear(e.target.value) }}/></p>
    <p>Til: <input type='number' value={selectedToYear}
      onChange={(e) => { setSelectedToYear(e.target.value) }}/></p>
    <Link to={`/s/${station.id}/from/${selectedFromYear}/to/${selectedToYear}`} >Oppdater</Link>
    <h2>Månadlege data</h2>
    <p>i perioden {stats.fromYear}-{stats.toYear}:</p>
    {table}
  </div>;
};
