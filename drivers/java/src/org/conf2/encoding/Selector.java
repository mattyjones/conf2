package org.conf2.encoding;

import java.net.MalformedURLException;
import java.net.URL;

/**
 *
 */
public class Selector {
    private String[] elements;
    private String query;
    private URL url;

    public Selector(String path) throws MalformedURLException {
        int paramsStart = path.indexOf('?');
        if (paramsStart > 0) {
            query = path.substring(0, paramsStart);
        } else {
            query = path;
        }
        elements = query.split("/");
    }

    public String[] getElements() {
        return elements;
    }
}
