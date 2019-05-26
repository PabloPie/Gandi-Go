## Disk
	- Est-ce vraiment utile de proposer un countDisk() ?
		=> pour savoir le nombre de pages qu'il faut query
	- Migration d'un disque d'un DC à l'autre ? (en bonus ?)

# Travail
	- 1ere partie IP+DISK+VM -> Tests
	- 2eme partie SSH+VLAN -> Tests
	- 3eme partie Terraform -> Tests
	- 4eme partie Doc 

# Bonus
	- Vlan
	- Disk snapshots
	- Migrations

# Problematiques
	- Shared structures for v4 and v5(= abstraction of the underlying objects to create a common representation)
	- IDs ( uuid in v5 vs int in v4 ) type conversion problems (= using strings for almost everything)
	- sync vs async (= waiting for creation operations to end)
	       async => Pointer vs value receiver(= immutability and concurrency problems)
		   try to stay close to api status => no status decoupling
	- Interfaces for unit testing (+ shared client interface for v4 and v5?)
                -> need a wrapper to mock api calls
