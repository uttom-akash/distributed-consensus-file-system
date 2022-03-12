import "./App.css";
import { useState } from "react";
import ChainTree from "./ChainTree";
import { Form, Input, Row, Col, Divider, Descriptions } from "antd";
import { Typography } from "antd";
import { Layout } from "antd";

function App() {
  const [blockhashmapper, setblockhash] = useState(null);
  const [selectedBlockHash, setselectedBlockHash] = useState(null);
  const [blocktree, setblocktree] = useState(null);

  const onSumbit = (e) => {
    e.preventDefault();
  };

  const onChange = (e) => {
    let json = JSON.parse(e.currentTarget.value);
    setblockhash(json.block_hash_mapper);
    setblocktree(json.block_tree);
  };

  const OnSelectBlock = (blockhash) => {
    setselectedBlockHash(blockhash);
    console.log(blockhash);
  };

  return (
    <Layout>
      <Row>
        <Col span={24}>
          <Typography.Title>Chain Visualizer</Typography.Title>
        </Col>
      </Row>
      <Divider></Divider>
      <Row justify="space-around" style={{ padding: "16px" }}>
        <Col span={7}>
          {!!selectedBlockHash && !!blockhashmapper[selectedBlockHash] && (
            <Row>
              <Col>
                <Descriptions title="Block Information" bordered column={1}>
                  <Descriptions.Item label="SerialNo">
                    {blockhashmapper[selectedBlockHash]["SerialNo"]}
                  </Descriptions.Item>
                  <Descriptions.Item label="BlockHash">
                    {selectedBlockHash}
                  </Descriptions.Item>
                  <Descriptions.Item label="PrevHash">
                    {blockhashmapper[selectedBlockHash]["PrevHash"]}
                  </Descriptions.Item>
                  <Descriptions.Item label="MinerID">
                    {blockhashmapper[selectedBlockHash]["MinerID"]}
                  </Descriptions.Item>
                  <Descriptions.Item label="TimeStamp">
                    {blockhashmapper[selectedBlockHash]["TimeStamp"]}
                  </Descriptions.Item>
                  <Descriptions.Item label="Nonce">
                    {blockhashmapper[selectedBlockHash]["Nonce"]}
                  </Descriptions.Item>
                  {!!blockhashmapper[selectedBlockHash]["Operations"] &&
                    blockhashmapper[selectedBlockHash]["Operations"].map(
                      (x, index) => (
                        <Descriptions.Item label={`Operation ${index}`}>
                          OperationId : {x["OperationId"]}
                          <br />
                          FileName : {x["FileName"]}
                          <br />
                          OperationType : {x["OperationType"]}
                          <br />
                          MinerID : {x["MinerID"]}
                          <br />
                          TimeStamp : {x["TimeStamp"]}
                          <br />
                          State : {x["State"]}
                          <br />
                          Record :{" "}
                          {String.fromCharCode(...x["Record"]).split("\x00")[0]}
                        </Descriptions.Item>
                      )
                    )}
                </Descriptions>
              </Col>
            </Row>
          )}
        </Col>
        <Divider orientation="left" type="vertical"></Divider>
        <Col span={16}>
          <Divider orientation="left">Json Data</Divider>
          <Row>
            <Col span={20}>
              <Form onSubmit={onSumbit}>
                <Form.Item
                  name="intro"
                  rules={[{ required: true, message: "Please input Intro" }]}
                >
                  <Input.TextArea showCount rows={10} onChange={onChange} />
                </Form.Item>

                {/* <Form.Item>
                <Button type="primary" htmlType="submit">
                  Show
                </Button>
              </Form.Item> */}
              </Form>
            </Col>
          </Row>
          <Divider orientation="left">Chain</Divider>
          <Row>
            <ChainTree blockTree={blocktree} OnSelectBlock={OnSelectBlock} />
          </Row>
        </Col>
      </Row>
    </Layout>
  );
}

export default App;
