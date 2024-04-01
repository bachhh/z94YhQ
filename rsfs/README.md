
# File system

A very basic file system, mostly modeling after ext4


## Basic file system layout

```
+----------------------------------------------------+
| SUPERBLOCK | I-MAP | D-Map | INODE BLOCK | DATA....|
+----------------------------------------------------+
```


## BLOCK layout

### Inode block

| Bytes | Name  | Note  |
| -     | -     | -     |
| Item1 | Item1 | Item1 |


### DATA block


## Feature list

### 1. ACL
 
- File ownership
- Access mode, Access Control list

