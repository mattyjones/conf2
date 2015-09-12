package org.conf2.schema.driver;

import org.conf2.schema.Handle;

/**
 *
 */
public class DriverHandle implements Handle {
    long handle;
    private Driver d;

    public DriverHandle(Driver d, long handle) {
        this.d = d;
        this.handle = handle;
    }

    public long getId() {
        return handle;
    }

    @Override
    public void Release() {
        d.releaseHandle(this);
    }
}
