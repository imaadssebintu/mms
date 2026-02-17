import psycopg2
import bcrypt
import uuid

def create_user(email, password, first_name, last_name, company_name, location, phone1, phone2):
    # Remote Database Connection Details
    db_config = {
        "host": "129.80.85.203",
        "port": 5432,
        "user": "imaad",
        "password": "Ertdfgxc",
        "dbname": "mms"
    }

    try:
        # Connect to the database
        conn = psycopg2.connect(**db_config)
        cur = conn.cursor()

        # Hash the password using bcrypt (matching the Go implementation cost of 14)
        salt = bcrypt.gensalt(14)
        hashed_password = bcrypt.hashpw(password.encode('utf-8'), salt).decode('utf-8')

        # Insert the user
        insert_query = """
        INSERT INTO users (email, password, first_name, last_name, company_name, location, phone1, phone2, updated_at)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s, NOW())
        RETURNING id;
        """
        
        cur.execute(insert_query, (
            email, 
            hashed_password, 
            first_name, 
            last_name, 
            company_name, 
            location, 
            phone1, 
            phone2
        ))
        
        user_id = cur.fetchone()[0]
        conn.commit()
        
        print(f"User created successfully with ID: {user_id}")
        
    except Exception as e:
        print(f"Error creating user: {e}")
    finally:
        if conn:
            cur.close()
            conn.close()

if __name__ == "__main__":
    # You can change these values as needed
    create_user(
        email="imaad.ssebintu@gmail.com",
        password="Ertdfgx@0",
        first_name="imaad",
        last_name="ssebintu",
        company_name="MMS Motors",
        location="Kampala",
        phone1="+256700752104",
        phone2=""
    )
