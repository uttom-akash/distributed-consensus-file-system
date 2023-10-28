import { Tree } from "antd";

const createNode = (hash) => ({
  title: hash,
  key: hash,
  children: [],
});

const createTreeData = (blockTree) => {
  console.log("createTreeData - In ", blockTree);

  if (!!!blockTree) return [];

  //Need to parse from genesis from block
  let genesisBlockHash = "5377672918f1bb6453dc3417e22966a0";

  let pblock = createNode(genesisBlockHash);
  let treeData = [pblock];
  let queue = [pblock];

  while (queue.length > 0) {
    let childrens = [];
    let currentBlock = queue[0];
    queue.shift();

    let sz = blockTree[currentBlock.key]?.length ?? 0;

    for (let index = 0; index < sz; index++) {
      let childBlockhash = blockTree[currentBlock.key][index];

      let childBlock = createNode(childBlockhash);

      childrens.push(childBlock);

      queue.push(childBlock);
    }
    currentBlock.children = childrens;
  }

  console.log("createTreeData - Out ", treeData);

  return treeData;
};

export default function ChainTree({ blockTree, OnSelectBlock }) {
  const onSelect = (selectedKeys, info) => {
    console.log("selected", selectedKeys, info);
    OnSelectBlock(selectedKeys[0]);
  };

  let treeData = createTreeData(blockTree);

  return (
    <Tree
      defaultExpandAll={true}
      showLine
      defaultExpandedKeys={["0-0-0"]}
      onSelect={onSelect}
      treeData={treeData}
    />
  );
}
