struct Inode {
    perm: u16,
    oid: u8,
    last_atime: u32, // last access time
    last_mtime: u32, // last modified time
    max_bloc: u64,
    i_block: [u32; 15],
}
