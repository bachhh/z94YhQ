
# File system

A very basic file system, mostly modeling after ext4

https://ext4.wiki.kernel.org/index.php/Ext4_Disk_Layout

TODO:

- ls
- cd
- mkdir
- touch
- write $FILE 
- cat $FILE 

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

A Data block 

### DIR BLOCK

Dir entry is a special type of data block, where the data content is a mapping from arbitary string to an INODE entry .

| Offset | Size | name     | Note                                      |
| -      | -    |          | -                                         |
| 0x0    | 32   | inod _no | Inode number of directory                 |
| 0      | 16   | entry_len|  total_len of entry |



| name |  | Note                                                                                   |
| -       | -    | -                                                                                      |
| current       | .    | The entry that point to Current directory                                              |
| 1       | ..   | The entry that point to Parent directory                                               |

| last    |      | End-of-block entry, with inode_number = 0 to signifies the end of directory data block |



## Feature list

- ACL, File ownership
- Access mode, Access Control list
- i_block 3-level indirect
- extent tree

