package org.conf2.yang;

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

    public static Meta findByIdent(MetaCollection c, String ident) {
        return findByIdent(new MetaCollectionIterator(c), ident);
    }

    public static boolean isLeaf(Meta m) {
        return m instanceof Leaf || m instanceof LeafList;
    }
}
