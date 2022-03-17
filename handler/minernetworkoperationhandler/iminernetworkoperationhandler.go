package minernetworkoperationhandler

type IMinerNetworkOperationHandler interface {
	DownloadChain()

	DisseminateOperations()

	DisseminateBlocks()
}
