import { darkTheme, GraphCanvas, lightTheme, type InternalGraphNode } from 'reagraph';
import './index.css';
import { useEffect, useState } from 'react';
import axios from 'axios';

// represents node table
interface PageNode {
  id: number;   
  url: string;
  title: string
}

// represent many-to-many relationship between pages
interface PageEdge {
  from_id: number; // match JSON keys exactly
  to_id: number;
}

// combines edges and nodes to create a graph object
interface LinkGraph {
  nodes: PageNode[];
  edges: PageEdge[];
}

type Props = {
  url: string;
};

const RelationshipGraph = ({ url }: Props) => {
  const [data, setData] = useState<LinkGraph | null>(null);

  useEffect(() => {
    const API_URL = `${import.meta.env.VITE_API_DOMAIN}/api/${import.meta.env.VITE_API_VER}/graph?domain=${url}`;
    console.log("graph URL:", API_URL);

    axios.get(API_URL)
      .then(res => setData(res.data))
      .catch(err => console.error(err));
  }, [url]);

  if (!data) return null;

  // map nodes for GraphCanvas
  const nodes = data.nodes.map(n => ({
    id: String(n.id),
    label: n.title,
  }));

  // map edges for GraphCanvas
  const edges = data.edges.map(e => ({
    id: `${e.from_id}->${e.to_id}`,
    source: String(e.from_id),
    target: String(e.to_id),
  }));

  const handleNodeClick = (node: InternalGraphNode) => {
    if (!data) return

    const res = data.nodes.find(o => o.id === Number(node.id))
    if (!res?.url) return

    window.open(res.url, '_blank')?.focus()
  }

  return (
    <div className="relationship-graph">
      <GraphCanvas 
        nodes={nodes}
        edges={edges}
        onNodeClick={(node: InternalGraphNode) =>handleNodeClick(node)}
        draggable
        theme={lightTheme}
      />
    </div>
  );
};

export default RelationshipGraph;
