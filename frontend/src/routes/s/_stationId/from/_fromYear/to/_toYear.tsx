import { useState } from 'react';
import {
  Link,
  LoaderFunctionArgs,
  useLoaderData,
  useRouteLoaderData
} from 'react-router-dom';
import { Station } from '../../../../..';

export const loader = ({params}: LoaderFunctionArgs) => {
  // fetch('api/params.stationId/from//to//');
  return {
    ...params,
    series: {
      'min': [-18, -17, -11.4, -7, -2, 0.5, 3.2, 4.2, -2, -14.4, -15, -18.4],
      'p50': [-3.2, -3.8, 0.8, 2.3, 11.5, 15.6, 14.2, 12, 9, 0.6, -3.6, -5],
      'max': [10.4, 8, 9.8, 14, 27.6, 28.6, 25.8, 20.5, 17.5, 18,5.5, 7.4]
    },
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
            <td>{stats.series.min[i]}</td>
            <td>{stats.series['p50'][i]}</td>
            <td>{stats.series.max[i]}</td>
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
