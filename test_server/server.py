# server.py
import asyncio
import logging
from asyncua import Server, ua
import random
from datetime import datetime


class OPCUATestServer:
    def __init__(self):
        self.server = Server()
        self._running = True

    async def init(self):
        await self.server.init()
        self.server.set_endpoint("opc.tcp://0.0.0.0:4840")
        self.server.set_security_policy([ua.SecurityPolicyType.NoSecurity])

        # # Set server name
        # await self.server.set_build_info(
        #     product_uri="urn:test:server",
        #     product_name="Test OPC UA Server",
        #     manufacturer_name="Test Manufacturer",
        #     software_version="1.0.0",
        # )

        # Get Objects node
        objects = self.server.nodes.objects

        # Create a custom namespace
        uri = "http://examples.test.com"
        idx = await self.server.register_namespace(uri)

        # Create nested structure
        # Level 1
        equipment = await objects.add_folder(idx, "Equipment")
        sensors = await objects.add_folder(idx, "Sensors")

        # Level 2 - Equipment
        pump = await equipment.add_folder(idx, "Pump1")
        motor = await equipment.add_folder(idx, "Motor1")

        # Level 2 - Sensors
        temp_sensors = await sensors.add_folder(idx, "Temperature")
        pressure_sensors = await sensors.add_folder(idx, "Pressure")

        # Add variables to pump
        self.pump_speed = await pump.add_variable(idx, "Speed", 0.0)
        self.pump_flow = await pump.add_variable(idx, "Flow", 0.0)
        await self.pump_speed.set_writable()
        await self.pump_flow.set_writable()

        # Add variables to motor
        self.motor_rpm = await motor.add_variable(idx, "RPM", 0.0)
        self.motor_temp = await motor.add_variable(idx, "Temperature", 0.0)
        await self.motor_rpm.set_writable()
        await self.motor_temp.set_writable()

        # Add temperature sensors
        self.temp1 = await temp_sensors.add_variable(idx, "Temp1", 0.0)
        self.temp2 = await temp_sensors.add_variable(idx, "Temp2", 0.0)
        await self.temp1.set_writable()
        await self.temp2.set_writable()

        # Add pressure sensors
        self.pressure1 = await pressure_sensors.add_variable(idx, "Pressure1", 0.0)
        self.pressure2 = await pressure_sensors.add_variable(idx, "Pressure2", 0.0)
        await self.pressure1.set_writable()
        await self.pressure2.set_writable()

        # Add timestamp
        self.timestamp = await objects.add_variable(idx, "CurrentTime", datetime.now())

    async def update_values(self):
        while self._running:
            try:
                # Update pump values
                await self.pump_speed.write_value(50 + random.uniform(-5, 5))
                await self.pump_flow.write_value(100 + random.uniform(-10, 10))

                # Update motor values
                await self.motor_rpm.write_value(1750 + random.uniform(-50, 50))
                await self.motor_temp.write_value(65 + random.uniform(-2, 2))

                # Update temperature sensors
                await self.temp1.write_value(25 + random.uniform(-1, 1))
                await self.temp2.write_value(28 + random.uniform(-1, 1))

                # Update pressure sensors
                await self.pressure1.write_value(2.5 + random.uniform(-0.1, 0.1))
                await self.pressure2.write_value(2.8 + random.uniform(-0.1, 0.1))

                # Update timestamp
                await self.timestamp.write_value(datetime.now())

                await asyncio.sleep(1)
            except Exception as e:
                logging.error(f"Error updating values: {e}")
                await asyncio.sleep(1)

    async def start(self):
        await self.init()
        async with self.server:
            update_task = asyncio.create_task(self.update_values())
            while self._running:
                await asyncio.sleep(1)
            update_task.cancel()

    def stop(self):
        self._running = False


async def main():
    server = OPCUATestServer()
    await server.start()


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    asyncio.run(main())
