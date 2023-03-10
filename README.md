# Usage

```
Usage: tailb <command>

Flags:
  -h, --help    Show context-sensitive help.

Commands:
  tail --start="2023-03-09 16:02:34" --end="2023-03-10 16:02:34" --nlines=-1 <loadbalancer>
    Tail logs of given loadbalancer. (default command)

  ls
    List loadbalancers.

Run "tailb <command> --help" for more information on a command.
```

# Example

```
tailb lb_name
```

```
tailb lb_name | jq 'select(.ElbStatusCode!="200")'
```

```
tailb lb_name -f TargetIP,ElbStatusCode
```