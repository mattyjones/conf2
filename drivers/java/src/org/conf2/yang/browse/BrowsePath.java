package org.conf2.yang.browse;

/**
 *
 */
public class BrowsePath {
    public Segments[] segments;
    public String Url;
    public BrowsePath(String p) {
        if (p == null || p.length() == 0) {
            segments = new Segments[0];
            return;
        }
        int qmark = p.indexOf('?');
        if (qmark >= 0) {
            Url = p.substring(0, qmark);
        } else {
            Url = p;
        }
        String[] segs = p.split("/");
        segments = new Segments[segs.length];
        for (int i = 0; i < segs.length; i++) {
            segments[i] = new Segments(this, segs[i]);
        }
    }
}

class Segments {
    Segments(BrowsePath path, String segment) {
        this.path = path;
        // TODO: take out keys
        this.ident = segment;
    }
    BrowsePath path;
    String ident;
    String[] keys;
}