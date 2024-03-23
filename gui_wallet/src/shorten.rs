pub(crate) fn shorten(s: String) -> String {
    let mut s = s;
    if s.len() > 10 {
        s.truncate(10);
        s.push_str("...");
    }
    s
}