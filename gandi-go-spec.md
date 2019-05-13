Design proposal for Gandi Hosting GO driver

## Disk image
	- Image identifiée par un int dans l'APIv4
	- une image = sur un DC
	- utilite de list images? normalement l'utilisateur devrait avoir une liste des images dispos sans devoir faire un appel sur l'api
## Disk
	- Est-ce vraiment utile de proposer un countDisk() ?
	- Il y a une histoire de noyaux dispo dans tel ou tel datacenter dans l'APIv4
	- Migration d'un disque d'un DC à l'autre ? (en bonus ?)


# Travail
	- 1ere partie IP+DISK+VM -> Tests
	- 2eme partie SSH+VLAN -> Tests
	- 3eme partie Terraform -> Tests
	- 4eme partie Doc 

# Problematiques
	- Shared structures for v4 and v5(= abstraction of the underlying objects to create a common representation)
	- IDs ( uuid in v5 vs int in v4 ) type conversion problems (= using strings for almost everything)
	- sync vs async (= waiting for creation operations to end)
	       async => Pointer vs value receiver(= immutability and concurrency problems)
	- Interfaces for unit testing (+ shared client interface for v4 and v5?)
                -> need a wrapper to mock api calls
