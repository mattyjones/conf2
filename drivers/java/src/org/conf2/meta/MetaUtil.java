package org.conf2.meta;

import java.util.Iterator;

/**
 *
 */
public class MetaUtil {
    public static Meta findByIdent(Iterator<Meta> i, String ident) {
        while (i.hasNext()) {
            Meta next = i.next();
            if (next.getIdent().equals(ident)) {
                return next;
            }
        }
        return null;
    }
}
