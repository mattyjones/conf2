package org.conf2.schema;

/**
 *
 */
public class SchemaError extends RuntimeException {
    public SchemaError(String msg) {
        super(msg);
    }
    public SchemaError(String msg, Throwable suberr) {
        super(msg, suberr);
    }
}
