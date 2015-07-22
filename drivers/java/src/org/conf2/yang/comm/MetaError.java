package org.conf2.yang.comm;

/**
 *
 */
public class MetaError extends RuntimeException {
    public MetaError(String msg) {
        super(msg);
    }
    public MetaError(String msg, Throwable suberr) {
        super(msg, suberr);
    }
}
