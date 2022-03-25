package peerclient

type IPeerClientHandler interface {
	DownloadChain()

	DisseminateOperations()

	DisseminateBlocks()
}
