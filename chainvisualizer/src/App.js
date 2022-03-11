import "./App.css";
import { Form } from "react-bootstrap";
import { Button } from "react-bootstrap";
import { useState } from "react";
import { Tree } from 'react-tree-graph';

const data ={
  "name": "Eve",
  "children": [
    {
      "name": "Cain"
    },
    {
      "name": "Seth",
      "children": [
        {
          "name": "Enos"
        },
        {
          "name": "Noam"
        }
      ]
    },
    {
      "name": "Abel"
    },
    {
      "name": "Awan",
      "children": [
        {
          "name": "Enoch"
        }
      ]
    },
    {
      "name": "Azura"
    }
  ]
}

function App() {
  const [chainJson, setchainJson] = useState("");

  const onSumbit = (e) => {
    e.preventDefault();
  };

  return (
    <div className="App">
      <header className="App-header">Chain Visualizer</header>
      <Form onSubmit={onSumbit}>
        <Form.Group className="mb-3" controlId="exampleForm.ControlTextarea1">
          <Form.Label>Chain Json</Form.Label>
          <Form.Control
            as="textarea"
            rows={3}
            onChange={(e) => setchainJson(e.currentTarget.value)}
          />
        </Form.Group>
        <Button variant="primary" type="submit">
          Submit
        </Button>
      </Form>
      <p>{chainJson}</p>
      <Tree data={data} height={1000} width={1000} />
      );
    </div>
  );
}

export default App;
