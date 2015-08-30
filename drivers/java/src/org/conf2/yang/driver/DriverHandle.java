package org.conf2.yang.driver;

import org.conf2.yang.Handle;

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

    @Override
    public void Release() {
        d.releaseHandle(this);
    }
}
