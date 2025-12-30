import { GraphCanvas } from 'reagraph';
import './index.css'

type Props = {}

const RelationshipGraph = (props: Props) => {
  return (
    <div className='relationship-graph'>
      <GraphCanvas
      nodes={[
        {
          id: 'n-1',
          label: '1'
        },
        {
          id: 'n-2',
          label: '2'
        }
      ]}
      edges={[
        {
          id: '1->2',
          source: 'n-1',
          target: 'n-2',
          label: 'Edge 1-2'
        }
      ]}
    />
    </div>
);
}

export default RelationshipGraph