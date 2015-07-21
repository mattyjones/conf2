package org.conf2.meta;

import java.util.Iterator;

/**
 *
 */
public class MetaCollectionIterator implements Iterator<Meta> {
    private Meta position;
    private boolean resolveProxies;
    private Iterator<Meta> currentProxy;

    public MetaCollectionIterator(MetaCollection c) {
        this(c, true);
    }

    public MetaCollectionIterator(MetaCollection c, boolean resolveProxies) {
        this.position = c.getFirstMeta();
        this.resolveProxies = resolveProxies;
    }

    @Override
    public boolean hasNext() {
        if (position != null) {
            return true;
        }
        if (currentProxy != null) {
            return currentProxy.hasNext();
        }
        return false;
    }

    @Override
    public Meta next() {
        while (position != null || currentProxy != null) {
            if (currentProxy != null) {
                if (currentProxy.hasNext()) {
                    return currentProxy.next();
                }
                currentProxy = null;
            } else {
                if (resolveProxies) {
                    if (position instanceof Proxy) {
                        currentProxy = ((Proxy) position).resolveProxy();
                        position = position.getSibling();
                    } else {
                        Meta next = position;
                        position = position.getSibling();
                        return next;
                    }
                } else {
                    Meta next = position;
                    position = position.getSibling();
                    return next;
                }
            }
        }
        return null;
    }
}
