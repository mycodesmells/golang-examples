# Bonus

## Go 1.8

Adding value to ENV causes duplicates for given key:

    $ USERNAME=Jane go run main.go 
    Jane
    Oscar

## Go 1.9 RC1

Go 1.9 fixed that:

    $ USERNAME=Jane go1.9rc1 run main.go 
    Oscar
